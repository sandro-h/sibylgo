import fetch from 'node-fetch';

export async function formatTodos(restUrl: string, text: string): Promise<string[]> {
    const res = await fetch(
        `${restUrl}/format`,
        {
            method: 'POST',
            headers: {'content-type' : 'text/plain'},
            body: Buffer.from(text).toString('base64')
        }
    );
    const body = await res.text();
    return body.split(/\r?\n/);
}

export async function foldTodos(restUrl: string, text: string): Promise<string[]> {
    const res = await fetch(
        `${restUrl}/fold`,
        {
            method: 'POST',
            headers: {'content-type' : 'text/plain'},
            body: Buffer.from(text).toString('base64')
        }
    );
    const body = await res.text();
    return body.split(/\r?\n/);
}

export async function preview(restUrl: string, text: string): Promise<string[]> {
    const res = await fetch(
        `${restUrl}/preview`,
        {
            method: 'POST',
            headers: {'content-type' : 'text/plain'},
            body: Buffer.from(text).toString('base64')
        }
    );
    const data = await res.json();
    return data;
}

export async function cleanTodos(restUrl: string): Promise<void> {
    await fetch(`${restUrl}/clean`, { method: 'POST' });
}

export async function trashTodos(restUrl: string): Promise<void> {
    await fetch(`${restUrl}/trash`, { method: 'POST' });
}
