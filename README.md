StaticWeb
=========

StaticWeb is a basic static files web server written in Go.

It automatically serves the content of the folder corresponding to the domain of the website.

Usage
-----

``./StaticWeb -addr="[Address to listen: 127.0.0.1]" -port="[Port to listen: 80]" -dir="[Absolute or relative path to the root directory to serve: .]" -log="[Absolute or relative path to the log file. Leave empty for stdout]"``

The program will serve the files in the folder (inside the given root directory) with the same name as the domain from which the request originated.

To serve multiple sites with the same content, you can use symbolic links.