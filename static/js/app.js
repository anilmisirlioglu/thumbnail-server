const templates = {
    'no-template': {
        url: null,
        y: null,
        selector: null,
        height: null,
        width: null,
        hide: null
    },
    github: {
        url: 'https://github.com',
        y: null,
        selector: null,
        width: 1440,
        height: 900,
        hide: [
            'body > div:first-child'
        ]
    }
}

const elements = ['url', 'width', 'height', 'y', 'selector', 'hide', 'template']

const [
    url, width, height, y,
    selector, hide, select
] = elements.map(element => document.getElementById(element))

select.onchange = () => {
    const template = templates[select.value] || templates['no-template']

    for(let [key, value] of Object.entries(template)){
        if(Array.isArray(value)) value = value.join('\n')

        window[key].value = value;
    }
}