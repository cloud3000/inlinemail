# inlinemail

Send email from the command line.

The best way to send emails from the command line.

# Introduction

InLineMail is a simple and efficient program to send emails.

## How it works

Two files must be edited:

- smtpconf.json -- To be configured for your smtp server.
- inlinemail.html -- The html contents to be in the email.

## At the command line

```bash
./inlinemail -f from@example.com -t to.someone@example.com -s "The subject line"
```
