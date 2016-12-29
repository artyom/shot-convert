Command shot-convert watches specified directory for new png screenshots and
converts them to jpeg images, leaving original screenshot intact.

Install:

	go get -u github.com/artyom/shot-convert

Create file `~/Library/LaunchAgents/org.acme.shot-convert.plist` with the
following content (adjust paths to match your environment):

	<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	<dict>
	    <key>Label</key>
	    <string>org.acme.shot-convert</string>
	    <key>Program</key>
	    <string>/Users/johndoe/go/bin/shot-convert</string>
	    <key>ProgramArguments</key>
	    <array>
			<string>shot-convert</string>
			<string>-dir=/path/to/screenshot-dir</string>
	    </array>
	    <key>KeepAlive</key>
	    <true/>
	</dict>
	</plist>

Run service:

	launchctl load ~/Library/LaunchAgents/org.acme.shot-convert.plist
