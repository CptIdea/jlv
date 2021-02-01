# JSON Log Viewer

CLI utility for viewing logs in JSON format

![img.png](https://sun9-33.userapi.com/impg/GsOJpyvsncmhE6CQN4jjtv8EtwvuUHQx6k45eQ/To4Daa-Seks.jpg?size=960x689&quality=96&proxy=1&sign=fa813f92fb3810e3437f30670d4b4e0e&type=album) ![img_1.png](https://sun9-47.userapi.com/impg/XZPtfbtkFjn-ISC7AARyl3n-CLM08o4p7hO6aQ/iAxHx9UiDjA.jpg?size=1280x472&quality=96&sign=3dadfe6f0f4c2b8815bb291643e366c7&type=album)

# Usage

run `jlv -F "PATH"` to simple view jsoned logs

`-t` short time

`-c` max count

`-s` show statistics

### `Filters(first char):`
`-w(arning)`
`-i(nfo)`
`-p(anic)`
`-e(rror)`
`-f(atal)`

`-I` invert filters



### `Sources:`
`-F(ile)`
`-T(CP)`
`TCP working with ALogger go package`

`-C` continuous reading of the source

# Features
### Tags
`time`
`prefix`
`level`
`msg`

### Levels
`warning`
`panic`
`fatal`
`error`
`info`

### Users info
If you want to add your own tag, it will be displayed as tag=data
