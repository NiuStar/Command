package Command

import (
	"io"
)

func ExecCommandWithAction(cmd,dir string,action func(string) string) error {
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
			r,err := child.ReadAll()
			if err != nil {
				if err == io.EOF {
					return
				}
				return
			}
			if action != nil {
				str := action(r)
				if len(str) > 0 {
					child.Send(str)
				}
			}
			//return
		}
	}()
	if err := child.Wait(); err != nil {
		return err
	}
	return nil
}

