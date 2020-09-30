# App

A sample 12 Facter Application.

## Usage

Download:

```
mkdir $GOPATH/src/github.com/kelseyhightower
cd $GOPATH/src/github.com/kelseyhightower 
git clone https://github.com/kelseyhightower/app.git
```

Generate TLS certificates:

```
$ go run certgen/main.go
```
```
wrote ca.pem
wrote ca-key.pem
wrote server.pem
wrote server-key.pem
```

### Build and Run

```
$ go build -o server ./monolith
```

```
$ ./server
```

```
2016/04/15 06:34:12 Starting server...
2016/04/15 06:34:12 HTTP service listening on 0.0.0.0:5000
2016/04/15 06:34:12 Health service listening on 0.0.0.0:5001
2016/04/15 06:34:12 Started successfully.
```

### Test with cURL

```
$ curl --cacert ./ca.pem -u user https://127.0.0.1:5000/login
```
```
Enter host password for user 'user':
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE0NjA5ODcxOTcsImlhdCI6MTQ2MDcyNzk5NywiaXNzIjoiYXV0aC5zZXJ2aWNlIiwic3ViIjoidXNlciJ9.x3oFhRhWk5CGYfGcrNctPGWCENEsXpUuKPDQU2ZOLCY
```

> type "password" at the prompt

```
curl --cacert ./ca.pem -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE0NjA5ODcxOTcsImlhdCI6MTQ2MDcyNzk5NywiaXNzIjoiYXV0aC5zZXJ2aWNlIiwic3ViIjoidXNlciJ9.x3oFhRhWk5CGYfGcrNctPGWCENEsXpUuKPDQU2ZOLCY' https://127.0.0.1:5000/
```
```
<h1>Hello</h1>
```
## ❤️ Contributors

<table>
  <tr>
    <td align="center"><a href="https://github.com/kelseyhightower"><img src="https://avatars0.githubusercontent.com/u/1123322?s=460&u=e9afe6d5b9bd6c20bbaf40f76b6188619fc24436&v=4" width="100px;" alt="Kelsey Hightower"/><br /><sub><b>Kelsey Hightower</b></sub></a><br /></td>
     </tr>
  </table>
