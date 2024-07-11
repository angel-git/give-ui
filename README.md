## GIVE-UI

Simple UI for the `SPT GIVE` command.

### How it works

1. With the server and preferably with Tarkov running, open the app
2. Use the form to connect to the server and select your character
3. Select the item you want to receive. The quantity is always set to the maximum stack size
4. You will receive a message with the item/s

### Development

```shell
templ generate && wails dev
```

### Release

- Update version in `wails.json`
- Update version in `server-mod/src/mod.ts`
- commit and push (TODO: automate this in future)
- Create a new release with proper tag
- Github action will take over and upload the zip
