# Thumbnail Server

Website Thumbnail & Screenshot Generator written in **Go** with
using [Chrome DevTools Protocol](https://github.com/chromedp/chromedp)


## `/screenshot` Endpoint

| Parameter |  Type  | Default | Description |
| --------- | ------ | ------- | ----------- |
| url       | string | -       | URL of the website that will generate the screenshot | 
| width     | int    | 1920    | Width of the screenshot |
| height    | int    | 1080    | Height of the screenshot |
| quality   | int    | 100     | Quality of the screenshot | 
| bgColor   | string | #000000 | Background color if available (If value is #000000 (default) means the no background) | 
| scrollY   | int    | -       | Is the pixel along the vertical axis (Y) |
| hide      | string | -       | Selectors of elements to be hidden |

## Preview

![image](https://cdn.discordapp.com/attachments/771673727473156109/808857787874279454/unknown.png)
