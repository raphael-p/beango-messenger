package client

var LoginPage string = `<span class="logo"><span>> Beango Messenger </span></span>
	<div id="login-form">
		<form hx-ext="json-enc">
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
		<div id=chat-list class="homepage-column">
			{{ range .Chats }}
			{{ block "chat-list" .}}
				<div 
					class="chat-selector list-item"
					hx-get="/home/chat/{{ .ID }}?name={{ .Name }}" 
					hx-target="#chat"
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
	<table id="message-table" class="homepage-column">` + messageRows + `</table>
	<div class="message-bar">
		<span class="message-prompt">> </span>
		<textarea
			class="message-input"
			placeholder="Type your message"
			hx-on:keypress="handleKeypress(event)"
			name="content"
			hx-post="/home/chat/{{ .ID }}/sendMessage"
			hx-trigger="send-message"
			hx-swap="none"
			hx-on::after-request="if(event.detail.successful) this.value = '';"
		></textarea>
	</div>` + newMessageFetcher

var MessagePaneRefresh string = newMessageFetcher + `
	<table id="message-table" hx-swap-oob="beforeend">` + messageRows + `</table>`

var messageRows string = `{{ range .Messages }}
	{{ block "message-list" .}}
		<tr class="list-item">
			<td class="cue">{{ .UserDisplayName }}</td>
			<td class="message">{{ .Content }}</td>
		</tr>
	{{ end }}
	{{ end }}`

var newMessageFetcher string = `<div 
	hx-get="/home/chat/{{ .ID }}?from={{ .FromMessageID }}&refresh"
	hx-swap="outerHTML"
	class="chat-selector list-item"
	hx-trigger="chat-refresh from:document, every 5s"
	/>`
