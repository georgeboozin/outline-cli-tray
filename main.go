package main

import (
    "fmt"
    "github.com/getlantern/systray"
    "log"
    "os"
    "os/exec"
    "syscall"
)

func startCmd(args []string) int {
    cmd := exec.Command("go", args...)

    cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    err := cmd.Start()
    if err != nil {
        log.Fatal(err)
    }

    return cmd.Process.Pid
}
func killProcess(pid int) {
    syscall.Kill(-pid, syscall.SIGINT)
}

func main() {
    accessKey := os.Args[1]
    outlineCli := "github.com/Jigsaw-Code/outline-sdk/x/examples/outline-cli@latest"
    extra := os.Args[2:]

    connectedIcon, err := os.ReadFile("./outline-connected-icon.svg")
    if err != nil {
        log.Fatal(err)
    }
    systray.SetIcon(connectedIcon)

    disconnectedIcon, err := os.ReadFile("./outline-disconnected-icon.svg")
    if err != nil {
        log.Fatal(err)
    }

    onReady := func() {
        args := []string{"run", outlineCli, "-transport", accessKey}
        args = append(args, extra...)
        pid := startCmd(args)

        go func() {
            isConneted := true
            statusItem := systray.AddMenuItem("Disconnect", "Connect or Disconnect actions")

            systray.AddSeparator()

            quitItem := systray.AddMenuItem("Quit", "Disconnect and Quit")

            for {
                select {
                case <-statusItem.ClickedCh:
                    if isConneted {
                        killProcess(pid)
                        statusItem.SetTitle("Connect")
                        isConneted = false
                        systray.SetIcon(disconnectedIcon)
                    } else {
                        pid = startCmd(args)
                        statusItem.SetTitle("Disconnect")
                        isConneted = true
                        systray.SetIcon(connectedIcon)
                    }
                case <-quitItem.ClickedCh:
                    fmt.Println("Requesting quit")
                    killProcess(pid)
                    systray.Quit()
                    fmt.Println("Finished quitting")
                    return
                }
            }
        }()
    }

    systray.Run(onReady, nil)
}
