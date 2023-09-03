package client

var LoginPage string = `<span class="logo"><span>> Beango Messenger </span></span>
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

var HomePage string = `<div id="chat-container">
	<div id="sidebar">
		<span class="heading-1">Chats</span>
		<div id=chat-list>
			{{ range .Chats }}
			{{ block "chat-list" .}}
				<div 
					hx-get=/home/chat/{{ .ID }}/{{ .Name }} 
					hx-target=#chat class="chat-selector list-item"
				>
					[{{ .Type}}] <b>{{ .Name }}</b>
				</div>
			{{ end }}
			{{ end }}
		</div>
	</div>
	<div id="chat"></div>
	</div>`

var MessagePane string = `<span class="heading-1">{{ .Name }}</span>
	<table>
		{{ range .Messages }}
		{{ block "message-list" .}}
			<tr class="list-item">
				<td class="cue">{{ .UserDisplayName }}</td>
				<td>{{ .Content }}</td>
			</tr>
		{{ end }}
		{{ end }}
	</table>
	<div class="message-bar">
		<span class="message-prompt">> </span>
		<textarea 
			class="message-input" 
			name="content" 
			hx-post=/home/chat/{{ .ID }}/sendMessage 
			hx-trigger="send-message"
			placeholder="Type your message"
		></textarea>
		<div 
			hx-get=/home/chat/{{ .ID }}/{{ .Name }} 
			hx-target=#chat class="chat-selector list-item"
			hx-trigger="chat-refresh from:document"
		/>
	</div>`
