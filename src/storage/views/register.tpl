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
        <h2><span class="fontawesome-lock"></span>Sign On</h2>
        <fieldset>
          <form action="/register" method="POST">
            <p><label for="email">E-mail address</label></p>
              <p><input type="email" name="email" placeholder="mail@address.com"></p>
              <p><label for="password">Password</label></p>
              <p><input type="password" name="password" placeholder="password"></p>
              <p><label for="password">Retype password</label></p>
              <p><input type="password" name="re-password" placeholder="password"></p>
              <p><input type="submit" value="Sign On"></p>
            </form>
        </fieldset>
      </div> 
    </div>
  </body> 
</html>