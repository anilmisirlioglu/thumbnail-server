const defaultSizeMetrics = {
    width: 1920,
    height: 1080
}

const templates = {
    'no-template': {
        url: null,
        y: null,
        selector: null,
        height: null,
        width: null,
        hide: null,
        bgColor: null,
        quality: null
    },
    github: {
        url: 'https://github.com',
        y: null,
        selector: null,
        width: 1920,
        height: 1080,
        hide: [
            'body > div:first-child'
        ],
        bgColor: null,
        quality: 100
    }
}

const elements = ['url', 'width', 'height', 'y', 'selector', 'hide', 'template', 'bgColor', 'quality']

const [
    url,
    width,
    height,
    y,
    selector,
    hide,
    select,
    bgColor,
    quality
] = elements.map(element => document.getElementById(element))

select.onchange = () => {
    const template = templates[select.value] || templates['no-template']

    for(let [key, value] of Object.entries(template)){
        if(Array.isArray(value)) value = value.join('\n')

        const element = window[key]
        element.value = value
        element.dispatchEvent(new Event('change'))
    }
}

let status = true
const prevMetrics = {
    width: null,
    height: null
}

bgColor.onchange = () => {
    const isReadyOnly = bgColor.value !== '#000000'
    width.readOnly = isReadyOnly
    height.readOnly = isReadyOnly
    if(isReadyOnly){
        for(const key of Object.keys(prevMetrics)){
            const curr = window[key]
            const value = curr.value
            if(value !== '' && status) prevMetrics[key] = value

            curr.value = defaultSizeMetrics[key]
        }

        status = false
    }else{
        for(const [key, value] of Object.entries(prevMetrics)){
            const curr = window[key]
            if(curr.value !== value){
                curr.value = value
            }
        }

        status = true
    }
}

for(const key of Object.keys(prevMetrics)){
    const element = window[key]
    element.onchange = () => {
        prevMetrics[key] = element.value
    }
}