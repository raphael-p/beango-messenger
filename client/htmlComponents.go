package client

var loginPage string = `<span class="title"><span>> Beango Messenger </span></span>
	<div id="login-form">
		<form hx-ext='json-enc'>
			<div class="form-row">
				<label for="username">Username:</label>
				<input type="text" name="username">
			</div>
			<div class="form-row">
				<label for="password">Password:</label>
				<input type="password" name="password">
			</div>
			<div class="form-row button-row">
				<button hx-post="/login/login" type="submit" hx-swap="none">Log In</button>
				<button hx-post="/login/signup" type="submit" hx-swap="none">Sign Up</button>
			</div>
			<div id="errors" class="error"></div>
		</form>
	</div>`

var homePage string = `<h1>You are home!</h1>`
