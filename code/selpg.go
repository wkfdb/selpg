package main

import (
	"github.com/spf13/pflag"
	"fmt"
	"io"
	"os"
	"bufio"
	"os/exec"
)

//import flag "github.com/ogier/pflag"

type selpg_args struct {
	s_page int
	e_page int
	len_page int
	page_type bool
	dest string
	filename string
}

var program_name string
var default_len int = 10

func main() {
	program_name = os.Args[0]
	var args selpg_args
	Init(&args)
	Check_error(&args)
	//fmt.Println(args)
	Readfile(&args)
}

func Init(args *selpg_args) {
	pflag.IntVarP(&args.s_page,"start_page", "s", -1, "start_page.")
	pflag.IntVarP(&args.e_page,"end_page", "e", -1, "end_page.")
	pflag.IntVarP(&args.len_page,"Lines_pre_page", "l", default_len, "Lines per page.")
	pflag.BoolVarP(&args.page_type,"page_type", "f", false, "Page type")
	//-f 表示使用换页符来决定一页多大，-f与-l是互斥的
	pflag.StringVarP(&args.dest,"Print_destination", "d", "", "Print destination")
    pflag.Usage = Usage
	pflag.Parse()
}

func Check_error(args *selpg_args) {

	if args.s_page == -1 || args.e_page == -1 {
		os.Stderr.Write([]byte("Must input -s -e\n"))
    	pflag.Usage()
		os.Exit(0)
	}
	if args.s_page < 1 {
		fmt.Fprintf(os.Stderr,"invalid start page\n")
    	pflag.Usage()
    	os.Exit(1)
	}
	
	if 	args.e_page < 1 {
		fmt.Fprintf(os.Stderr,"invalid end page\n")
    	pflag.Usage()
    	os.Exit(2)
	}

	if args.len_page < 1 {
		fmt.Fprintf(os.Stderr,"invalid page len\n")
    	pflag.Usage()
    	os.Exit(3)
	}
	if args.s_page > args.e_page {
    	fmt.Fprintf(os.Stderr,"start page larger than end page\n")
    	pflag.Usage()
    	os.Exit(4)
  	}

 	if args.len_page != default_len && args.page_type {
    	fmt.Fprintf(os.Stderr,"%s:-f and -l can not input together\n",program_name)
		pflag.Usage()
		os.Exit(5)
	}

 	if len(pflag.Args()) > 1 {
    	fmt.Fprintf(os.Stderr,"%s:too much args\n",program_name)
    	pflag.Usage()
    	os.Exit(6)
  	}

}

func Readfile(args *selpg_args) {

	var stdin io.WriteCloser
	var err error
	var cmd *exec.Cmd

	//表示输出到哪里哪里
    if args.dest != "" {
		cmd = exec.Command("cat", "-n")
		stdin, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		stdin = nil
	}

	//如果有输入文件名就读取文件
    if pflag.NArg() > 0 {
	    args.filename = pflag.Arg(0)
	    input,err := os.Open(args.filename)

	    if err != nil {
	      fmt.Println(err)
	      os.Exit(7)
	    }

	    bfreader := bufio.NewReader(input)
	    //如果使用换页符来决定每一页
	    if args.page_type {
	    	page_num := 1;
	        for page_num <= args.e_page {
	        	//获取一页内容
	        	line_string, err := bfreader.ReadString('\f')
	        	//line_string, err := bfreader.ReadString('\n')
	        	if err != nil && err != io.EOF {
	          		fmt.Println(err)
	          		os.Exit(8)
	        	}

	        	if err == io.EOF {
	          		break
	        	}
	        	//若页数处于读取的范围内，写入目标文件或屏幕
	        	if page_num >= args.s_page {
	          		if args.dest != "" {
	            		stdin.Write([]byte(string(line_string) + "\n"))
		        	} else {
			        	fmt.Print(string(line_string))
		        	}
	        	}
	        	page_num ++
	      	}
	    } else {
	    	//一行一行地读
			line_num := 0
	      	page_num := 1

			for {
				//读取一行内容
				line_string, _, err := bfreader.ReadLine()
				if err != nil && err != io.EOF {
					fmt.Println(err)
					os.Exit(9)
				}
				if err == io.EOF {
					break
				}
	        	if page_num >= args.s_page && page_num <= args.e_page {
	        		//写入目标文件或屏幕
	          		if args.dest != "" {
	            		stdin.Write([]byte(string(line_string) + "\n"))
	          		} else {
	            		fmt.Println(string(line_string))
	          		}
	        	}
	        	line_num ++
	        	if line_num == args.len_page {
	          		page_num ++
	          		line_num = 0
	        	}
	        	if page_num > args.e_page {
	          		break
	        	}
			}
		}
	} else {
		//否则读取标准输入
		bfscanner := bufio.NewScanner(os.Stdin)
		line_num := 0
    	page_num := 1
		out_string := ""
		for bfscanner.Scan() {
			//读取标准输入内容
			line_string := bfscanner.Text()
			line_string += "\n"
      		if page_num >= args.s_page && page_num <= args.e_page {
        		out_string += line_string
      		}
      		line_num ++
      		if line_num == args.len_page {
        		page_num ++
        		line_num = 0
      		}
		}
		//写入目标文件或屏幕
		if args.dest != "" {
      		stdin.Write([]byte(string(out_string) + "\n"))
    	} else {
      		fmt.Println(string(out_string))
    	}
	}

	if args.dest != "" {
		stdin.Close()
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func Usage() {
  fmt.Fprintf(os.Stderr,
		`usage: [-s start_page page(>=1)] [-e end_page page(>=s)] [-l length of page(default 72)] [-f type of file(default 1)] [-d dest] [filename specify input file]
`)
}