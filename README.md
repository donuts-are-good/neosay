
![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)

# 🎉 Neosay
Neosay is an easy CLI tool that allows you to send messages and code from files or stdin to Matrix.

## 📥 Click the banner to download

[![image](https://user-images.githubusercontent.com/96031819/230751959-ee18e3bf-9559-4aa7-a6f4-da5faeabac5f.png)](https://github.com/donuts-are-good/neosay/releases/latest) 

Download a copy of Neosay for your system by clicking the big green banner above.


## 🚀 Building It Yourself
To build Neosay, you'll need to clone the repository and build the binary:

```
git clone https://github.com/donuts-are-good/neosay.git
cd neosay
go build
```
## 📝 Configuration (Environment Variables)
You may want to use Neosay with environment variables instead of a config file. 

```
echo "Hello, world!" | ./neosay
```

You can load in your configuration details as environment variables like this in your environment:

```
export MATRIX_HOMESERVER_URL="https://matrix.example.com"
export MATRIX_USER_ID="@yourusername:example.com"
export MATRIX_ACCESS_TOKEN="your_access_token"
export MATRIX_ROOM_ID="!yourroomid:example.com"
```


## 📝 Configuration (JSON)
Neosay can optionally use a JSON configuration file to store the Matrix homeserver URL, user ID, access token, and room ID. This makes it easy to switch between rooms. 

**Note:** Be careful using this method, anybody can inmpersonate you with this information. It's best to make a different account for this.

Here's an example of the JSON configuration file:

```
{
  "homeserverURL": "https://matrix.example.com",
  "userID": "@yourusername:example.com",
  "accessToken": "your_access_token",
  "roomID": "!yourroomid:example.com"
}
```
## 🌟 Examples

### Sending a simple message
`echo "Hello, world!" | ./neosay config.json`

### Sending a multi-line message
`cat yourfile.txt | ./neosay config.json`

### Sending a code block
`echo "Hello, world!" | ./neosay -c config.json`

### Sending a multi-line code block
`cat yourcodefile.py | ./neosay --code config.json`

## 🤖 Neosay in the wild
Here are some other ways to use Neosay:

- Share code snippets with your team in a Matrix room
- Send logs or debugging information from a script
- Share inspirational quotes or jokes
- Send notifications from your CI/CD pipelines or monitoring tools


## 📡 Contributing
I'll take help however you can get it to me, whether that be a pull request, a .patch file, a sticky note, smoke signals, bribes, you name it.

## 🖇️ License
`neosay` is licensed under the MIT license. Check out the LICENSE file for more info. If you dont get it don't worry about it.