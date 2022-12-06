package apikey

func SerializeRequest(method, host, path, body string) string {
	return method + "\n" + host + "\n" + path + "\n" + body
}
