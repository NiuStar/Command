package Command

import (
	"bytes"
	"fmt"
	"io"
	//"strings"
)
//没有交互的
func ExecCommandNoAction(cmd,dir string) error {
	fmt.Println("order:",cmd)
	child, err := CommandCreate(cmd)
	if err != nil {
		return err
	}
	child.Cmd.Dir = dir
	err = child.Start()
	if err != nil {
		return err
	}

	go func() {
		for {
			r,err :=child.ReadAll()
			if err == io.EOF {
				fmt.Println("Over")
				return
			} else if err != nil {
				fmt.Println("cmd:",cmd,"error:",err)
				return
			}
			fmt.Println(r)
		}
	}()
	if err := child.Wait(); err != nil {
		return err
	}
	return nil
}


func ExecCommandWithResult(cmd,dir string) (string,error) {
	fmt.Println("order:",cmd)
	child, err := CommandCreate(cmd)
	if err != nil {
		return "",err
	}
	child.Cmd.Dir = dir
	err = child.Start()
	if err != nil {
		return "",err
	}

	msg := make(chan string)
	go func() {
		msg1 := bytes.Buffer{}
		for {
			r,err :=child.ReadAll()
			if err == io.EOF {
				//fmt.Println("Over",msg1.String())
				msg <- msg1.String()
				return
			} else if err != nil {
				fmt.Println("cmd:",cmd,"error:",err)
				msg <- msg1.String()
				return
			}
			//fmt.Println("r:",r)
			msg1.WriteString(r)
		}
	}()
	if err := child.Wait(); err != nil {
		return <-msg,err
	}

	return <-msg,nil
}
