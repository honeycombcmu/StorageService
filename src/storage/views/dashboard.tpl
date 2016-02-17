<!DOCTYPE html>
<head>
	<meta charset="utf-8">
	<title>Dashboard</title>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
	<link rel="stylesheet" href="http://fonts.googleapis.com/css?family=Varela+Round">
    <link rel="stylesheet" type="text/css" href="/static/css/login.css">
</head>

<body>
	<h4>Welcome to honeycomb, {{.User_name}}</h4>
    <h4>{{.todo}}</h4>
    </div>
    	<form method="post" name="submit" enctype="multipart/form-data">
  			<input type="file" name="fileField"><br /><br />
  			<input type="submit" name="submit" value="Submit">
		</form>
    </div>
</body>