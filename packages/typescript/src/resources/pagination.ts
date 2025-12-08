import type { ListResponse } from '../types/index.js';

/**
 * Async iterator for auto-pagination of resources.
 * Yields individual items from pages, fetching new pages as needed.
 */
export class Paginator<T> implements AsyncIterableIterator<T> {
    private buffer: T[] = [];
    private nextCursor?: string;
    private hasMore = true;

    constructor(
        private readonly fetchNext: (cursor?: string) => Promise<ListResponse<T>>,
        firstPage?: ListResponse<T>
    ) {
        if (firstPage) {
            this.buffer = [...firstPage.data];
            this.hasMore = firstPage.has_more;
            this.nextCursor = firstPage.next_cursor;
        }
    }

    [Symbol.asyncIterator](): AsyncIterableIterator<T> {
        return this;
    }

    async next(): Promise<IteratorResult<T>> {
        if (this.buffer.length > 0) {
            return { done: false, value: this.buffer.shift()! };
        }

        if (!this.hasMore) {
            return { done: true, value: undefined };
        }

        const page = await this.fetchNext(this.nextCursor);
        this.buffer = [...page.data];
        this.hasMore = page.has_more;
        this.nextCursor = page.next_cursor;

        if (this.buffer.length === 0) {
            return { done: true, value: undefined };
        }

        return { done: false, value: this.buffer.shift()! };
    }
}
