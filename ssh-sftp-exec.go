package main

import (
  "fmt"
  "log"
  "os"
  "path"
  "time"
  "net"
  "flag"

  "github.com/pkg/sftp"

  "golang.org/x/crypto/ssh"
  "golang.org/x/crypto/ssh/terminal"
)

var (
   op = flag.String("op", "put", "operator")
   hostip = flag.String("host", "127.0.0.1", "remote host ip")
   username = flag.String("username", "test", "remote user name")
   userpasswd =  flag.String("userpasswd", "20190823", "remote user passwd")
   srcf =  flag.String("srcf", "./t.txt", "local path file name")
   tard = flag.String("tard", "./tmp/", "remote path")
   cmd = flag.String("cmd", "pwd", "remote path")
)

type SftpCli struct{
  sftpClient *sftp.Client
  sshSession *ssh.Session
}

func main(){
  flag.Parse()
  sc :=NewSftpCli()

  if *op == "put" {
     sc.put()
  }else if  *op == "get" {
     sc.get()
  }else if *op == "cmd" {
     sc.execcmd()
  }else{
     sc.terminalcmd()
  }

  defer sc.sftpClient.Close()
  defer sc.sshSession.Close()
}

func NewSftpCli() *SftpCli{
  sshSession,sftpClient, err := connect(*username, *userpasswd, *hostip, 22)
  if err != nil {
    log.Fatal(err)
  }
  return &SftpCli{
     sftpClient: sftpClient,
     sshSession: sshSession,
  }
}

func (s *SftpCli) put() error{

  srcFile, err := os.Open(*srcf)
  if err != nil {
    log.Fatal(err)
  }
  defer srcFile.Close()
  var remoteFileName = path.Base(*srcf)
  dstFile, err := s.sftpClient.Create(path.Join(*tard, remoteFileName))
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer dstFile.Close()
  buf := make([]byte, 1024)
  for {
    n, _ := srcFile.Read(buf)
    if n == 0 {
      break
    }
    dstFile.Write(buf)
  }

  fmt.Println("copy file to remote server finished!")
  return nil
}

func (s *SftpCli) get() error{

  srcFile, err := s.sftpClient.Open(*srcf)
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer srcFile.Close()

  var localFileName = path.Base(*srcf)
  dstFile, err := os.Create(path.Join(*tard, localFileName))
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer dstFile.Close()

  if _, err = srcFile.WriteTo(dstFile); err != nil {
    log.Fatal(err)
    return err
  }

  fmt.Println("copy file from remote server finished!")
  return nil
}

func connect(user, password, host string, port int) (*ssh.Session, *sftp.Client, error) {
  var (
    auth         []ssh.AuthMethod
    addr         string
    clientConfig *ssh.ClientConfig
    sshClient    *ssh.Client
    sftpClient   *sftp.Client
    sshSession   *ssh.Session
    err          error
  )
  // get auth method
  auth = make([]ssh.AuthMethod, 0)
  auth = append(auth, ssh.Password(password))

  clientConfig = &ssh.ClientConfig{
    User:    user,
    Auth:    auth,
    Timeout: 30 * time.Second,
    HostKeyCallback:func(hostname string, remote net.Addr, key ssh.PublicKey) error {
        return nil
    },
  }

  // connet to ssh
  addr = fmt.Sprintf("%s:%d", host, port)

  if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
    return nil,nil, err
  }

  // create sftp client
  if sftpClient, err = sftp.NewClient(sshClient); err != nil {
    return nil,nil, err
  }
  // create session
  if sshSession, err = sshClient.NewSession(); err != nil {
    return nil,nil, err
  }

  return sshSession,sftpClient, nil

}

func (s *SftpCli) execcmd() error {
  s.sshSession.Stdout = os.Stdout
  s.sshSession.Stderr = os.Stderr
  s.sshSession.Run(*cmd)
  return nil
}

func (s *SftpCli) terminalcmd(){
  fd := int(os.Stdin.Fd())
  oldState, err := terminal.MakeRaw(fd)
  if err != nil {
    panic(err)
  }
  defer terminal.Restore(fd, oldState)

  // excute command
  s.sshSession.Stdout = os.Stdout
  s.sshSession.Stderr = os.Stderr
  s.sshSession.Stdin = os.Stdin

  termWidth, termHeight, err := terminal.GetSize(fd)
  if err != nil {
    panic(err)
  }

  // Set up terminal modes
  modes := ssh.TerminalModes{
    ssh.ECHO:          1,     // enable echoing
    ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
    ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
  }

  // Request pseudo terminal
  if err := s.sshSession.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
    log.Fatal(err)
  }

  s.sshSession.Run(*cmd)

}
