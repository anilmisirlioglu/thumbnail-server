# Thumbnail Server

[![shields](https://img.shields.io/badge/made%20with-go-blue?logo=go&style=for-the-badge&logoColor=white)](https://golang.org)

Website Thumbnail & Screenshot Generator written in **Go** with
using [Chrome DevTools Protocol](https://github.com/chromedp/chromedp).

<p align="center">
  <img src="https://cdn.discordapp.com/attachments/771673727473156109/808857787874279454/unknown.png" width="738">
</p>

> [Click to try a live demo.](https://tn.anil.tools)
> NOTE: Fonts may appear bold due to heroku chrome buildpack.

## Features

- Generating the Website Screenshots & Thumbnails.
- Set the custom width, height and quality values of the screenshot.
- Ability to scroll the page on the Y axis.
- Set the background and customize the background color of the screenshot.
- Delete the elements you do not want from the screenshot with the sensitive selector.

## TODO

- [ ] Add the Chrome UI to images with backgrounds.
- [ ] Fullscreen page screenshot.
- [ ] Allowing custom height and width adjustments for images with backgrounds.

## Query Parameters

> for `/screenshot` Endpoint

| Parameter |  Type  | Default | Description |
| --------- | ------ | ------- | ----------- |
| url       | string | -       | URL of the website that will generate the screenshot | 
| width     | int    | 1920    | Width of the screenshot |
| height    | int    | 1080    | Height of the screenshot |
| quality   | int    | 100     | Quality of the screenshot | 
| bgColor   | string | #000000 | Background color if available (If value is #000000 (default) means the no background) | 
| scrollY   | int    | -       | Is the pixel along the vertical axis (Y) |
| hide      | string | -       | Selectors of elements to be hidden |


## Installation & Test

1. First, install repository
    ```sh
    git clone https://github.com/anilmisirlioglu/thumbnail-server && cd thumbnail-server
    ```

2. Then, run it.
    ```sh
    go run .
    ``` 

3. Finally, enter localhost address from your browser.