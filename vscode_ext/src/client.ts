import * as request from 'request';

export function formatTodos(restUrl: string, text: string): Promise<string[]> {
    return new Promise((resolve, reject) => {
        request.post(
            {
                headers: {'content-type' : 'text/plain'},
                url:     `${restUrl}/format`,
                body:    Buffer.from(text).toString('base64')
              },
              (error, response, body: string) => {
                if (error || response.statusCode !== 200) {
                    reject(`HTTP ${response.statusCode} ${error}`);
                }
                else if (!body) {
                    reject('No response body');
                }
                else {
                    resolve(body.split(/\r?\n/));
                }
            }
        );
    });
}

export function foldTodos(restUrl: string, text: string): Promise<string[]> {
    return new Promise((resolve, reject) => {
        request.post(
            {
                headers: {'content-type' : 'text/plain'},
                url:     `${restUrl}/folding`,
                body:    Buffer.from(text).toString('base64')
              },
              (error, response, body: string) => {
                if (error || response.statusCode !== 200) {
                    reject(`HTTP ${response.statusCode} ${error}`);
                }
                else if (!body) {
                    reject('No response body');
                }
                else {
                    resolve(body.split(/\r?\n/));
                }
            }
        );
    });
}

export function cleanTodos(restUrl: string): Promise<void> {
    return new Promise((resolve, reject) => {
        request.post(
            `${restUrl}/clean`,
            (error, response) => {
                if (error || response.statusCode !== 200) {
                    reject(`HTTP ${response.statusCode} ${error}`);
                }
                else {
                    resolve();
                }
            }
        );
    });
}

export function trashTodos(restUrl: string): Promise<void> {
    return new Promise((resolve, reject) => {
        request.post(
            `${restUrl}/trash`,
            (error, response) => {
                if (error || response.statusCode !== 200) {
                    reject(`HTTP ${response.statusCode} ${error}`);
                }
                else {
                    resolve();
                }
            }
        );
    });
}
