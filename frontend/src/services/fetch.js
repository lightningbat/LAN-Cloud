export default async function _fetch(route, body) {
    try {
        const response = await fetch(route, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(body)
        });
        if (response.status !== 200) {
            throw new Error(response.statusText);
        } else {
            
        }
    } catch (error) {
        console.log(error);
    }
}