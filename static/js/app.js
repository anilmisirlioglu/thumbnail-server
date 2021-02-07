const templates = {
    'no-template': {
        url: null,
        y: null,
        selector: null,
        height: null,
        width: null
    },
    github: {
        url: 'https://github.com',
        y: 0,
        selector: null,
        width: 1440,
        height: 900
    }
}

const elements = ['url', 'width', 'height', 'y', 'selector', 'template']

const [
    url, width, height,
    y, selector, select
] = elements.map(element => document.getElementById(element))

select.onchange = () => {
    const template = templates[select.value] || templates['no-template']

    for([key, value] of Object.entries(template)){
        window[key].value = value;
    }
}