<!--
<!DOCTYPE html>

<head>
	<meta charset="utf-8">
	<title>System Login</title>
	<link rel="stylesheet" type="text/css" href="/static/css/login.css">
</head>


<div class="loginpanel">
  <div class="txt">
    <input id="user" type="text" placeholder="Username" />
    <label for="user" class="entypo-user"></label>
  </div>
  <div class="txt">
    <input id="pwd" type="password" placeholder="Password" />
    <label for="pwd" class="entypo-lock"></label>
  </div>
  <div class="buttons">
    <input type="button" value="Login" />
    <span>
      <a href="javascript:void(0)" class="entypo-user-add register">Register</a>
    </span>
  </div>
  
  <div class="hr">
    <div></div>
    <div>OR</div>
    <div></div>
  </div>
</div>

<span class="resp-info"></span>

</body>
</html>
-->
<html lang="en-US">
  <head>
    <meta charset="utf-8">
    <title>Login</title>
    <link rel="stylesheet" href="http://fonts.googleapis.com/css?family=Varela+Round">
    <link rel="stylesheet" type="text/css" href="/static/css/login.css">
  </head>

  <body>
    <div class="container">
      <div id="login">
        <h2><span class="fontawesome-lock"></span>Sign In</h2>
        <fieldset>
        	<form action="/login" method="POST">
        		<p><label for="email">E-mail address</label></p>
            	<p><input type="email" name="email" placeholder="mail@address.com"></p>
            	<p><label for="password">Password</label></p>
            	<p><input type="password" name="password" placeholder="password"></p>
            	<p><input type="submit" value="Sign In"></p>
          	</form>
        	<form action="/register" method="GET">
        		<p><input type="submit" value="Register"></p>
        	</form>
        </fieldset>
      </div> 
    </div>
  </body>	
</html>
