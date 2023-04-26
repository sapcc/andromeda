const ENDPOINT = "https://endpoint-url-here.com"

class HTTPError extends Error {
    constructor(code, message) {
        super(message || code)
        this.name = "HTTPError"
        this.statusCode = code
    }
}

const encodeUrlParamsFromObject = (options) => {
    if (!options) return ""
    let encodedOptions = Object.keys(options)
        .filter((k) => options[k])
        .map((k) => `${encodeURIComponent(k)}=${encodeURIComponent(options[k])}`)
        .join("&")
    return `&${encodedOptions}`
}

// Check response status
const checkStatus = (response) => {
    if (response.status < 400) {
        return response
    } else {
        return response.json().then((message) => {
            const msg = message.error || message
            const error = new HTTPError(msg.code, msg.message)
            return Promise.reject(error)
        }, () => {
            const error = new HTTPError(response.status, response.statusText)
            error.statusCode = response.status
            return Promise.reject(error)
        })
    }
}

export const updateAttributes = (target, source) => {
    const res = {};
    Object.keys(target)
        .forEach(k => res[k] = (source[k] ?? target[k]));
    return res
}

export const nextPageParam = (lastPage, pages) => {
    return new URLSearchParams(
        lastPage.links?.find(x => {return x.rel === 'next'})?.href
    ).get('marker')
}

export const fetchAll = ({queryKey, pageParam, meta}) => {
    const [key, options] = queryKey
    // Support for useInfiniteQuery
    const query = encodeUrlParamsFromObject({...options, marker: pageParam})
    return fetch(`${meta.endpoint}/${key}?${query}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "X-Auth-Token": meta.token,
            Accept: "application/json",
        },
    })
        .then(checkStatus)
        .then((response) => {
            return response.json()
        })
}

export const fetchItem = ({queryKey, meta}) => {
    const [key, id, options] = queryKey
    const query = encodeUrlParamsFromObject(options)
    return fetch(`${meta.endpoint}/${key}/${id}?${query}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "X-Auth-Token": meta.token,
            Accept: "application/json",
        },
    })
        .then(checkStatus)
        .then((response) => {
            return response.json()
        })
}

export const updateItem = ({key, id, endpoint, formState, token}) => {
    // Converts a JavaScript value to a JavaScript Object Notation (JSON) string.
    const sendBody = JSON.stringify(formState)
    return fetch(`${endpoint}/${key}/${id}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
            "X-Auth-Token": token,
            Accept: "application/json",
        },
        body: sendBody,
    })
        .then(checkStatus)
        .then((response) => {
            return response.json()
        })
}

export const deleteItem = ({key, endpoint, id, token}) => {
    return fetch(`${endpoint}/${key}/${id}`, {
        method: "DELETE",
        headers: {
            "X-Auth-Token": token,
        },
    })
        .then(checkStatus)
}

export const createItem = ({key, endpoint, formState, token}) => {
    // Converts a JavaScript value to a JavaScript Object Notation (JSON) string.
    const sendBody = JSON.stringify(formState)
    return fetch(`${endpoint}/${key}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "X-Auth-Token": token,
            Accept: "application/json",
        },
        body: sendBody,
    })
        .then(checkStatus)
        .then((response) => {
            return response.json()
        })
}

export const login = ({endpoint, username, password, domain, project}) => {
    const identity = (username && password) ?
        {
            methods: ["password"],
            password: {
                user: {
                    name: username,
                    domain: {
                        name: domain,
                    },
                    password: password,
                }
            }
        } : {
            "methods": ["external"],
            "external": {}
        }

    const auth = {
        auth: {
            identity: identity,
            scope: {
                project: {
                    name: project,
                    domain: {
                        name: domain
                    }
                }
            }
        }
    }
    const sendBody = JSON.stringify(auth)
    return fetch(`${endpoint}/auth/tokens`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "X-User-Domain-Name": domain,
            Accept: "application/json",
        },
        body: sendBody,
    })
        .then(checkStatus)
        .then((response) => {
            return response.json().then(data => {
                return [response.headers.get("X-Subject-Token"), data.token]
            })
        })
}

