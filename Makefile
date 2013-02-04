README.html:	README.mkd
	pandoc -t html5 -s README.mkd > README.html
