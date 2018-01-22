![enter image description here](https://github.com/cheikhshift/bullscript/raw/master/bull.png)
# Bullscript
A command used to process .bs files into Web server source. The essence of this tool is to assemble a web application from other packages.

### How it works?
You fill an empty file with .bs calls and then use the command to convert the file into Go.
BS calls include HTTP request routing, run other commands during build and Go statements to run prior to server launch.

### BS calls

#### Index
Run a file with a bs call that has no parameters to explain its functionality.

1. [@onstart](#onstart)
2. [@i](#i)
3. [@prefix](#prefix)
4. [@path](#path)
5. [@listen](#listen)
6. [@listensecure](#listensecure)
7. [@redirect](#redirect)
8. [@run](#run)


### onstart
This bs call will run the specified Go statement on launch.
		
	Example : @onstart > gopkg.ExportedFunc("Hey I'm a string")

### i
This bs call will import the specified Go package. You may also use custom package names as in a  normal Go file.

	Example : @i > "path.com/gopackage/package"
	Example 2 : @i > customname "path.com/gopackage/package"

### prefix
This bs call will add the specified path prefix to any of the following paths. This allows for easy API grouping and loss of redundancy in specifying each path explicitly. The following example will add prefix `/home` to the only sub path. Use bs call `@end`to stop the command from adding the prefix to future paths. 

	Example: 
	@prefix > /home
		@path > /random/path > formPackage.Handler
	@end

The bs `@path` would be accessible with URL : `/home/random/path`

### path
This bs call will match the specified path with the specified handler.

	Example: @path > /random/path > form.Handler 

### listen
This bs call will specify which port the web application should listen and serve on.

	Example: @listen > 8080
### listensecure
This bs call will specify which port the web application should listen and serve on, as well as the path to your application's TLS files.

Syntax of bs call : `@listensecure > PORT >path_to_certificate_file > path_to_key_file`

	Example: @listensecure > 443 > server.cert > server.key

### redirect
This bs call will specify which port the specified port should redirect to.

	Example : @redirect > ORIGINAL_PORT > REDIRECT_PORT

### run
This bs call will run the specified terminal command.

	Example: @run > echo "Hello world!"

### Sample
Checkout the directory named `sample` for a full example of this command.

