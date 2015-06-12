// Type definitions for jQuery simplePagination.js v1.4
// Project: https://github.com/flaviusmatis/simplePagination.js
// Definitions by: Natan Vivo <https://github.com/nvivo/>
// Definitions: https://github.com/borisyankov/DefinitelyTyped

/// <reference path="../jquery/jquery.d.ts" />

interface SimplePaginationOptions {
    items?: number;
    itemsOnPage?: number;
    pages?: number;
    displayedPages?: number;
    edges?: number;
    currentPage?: number;
    hrefTextPrefix?: string;
    hrefTextSuffix?: string;
    prevText?: string;
    nextText?: string;
    cssStyle?: string;
    selectOnClick?: boolean;
    onPageClick?: (page: number, event: any) => void;
    onInit?: () => void;
}

interface JQuery {
    /**
     * The Vue Constructor
     * http://vuejs.org/api/index.html
     */
    pagination(options?: SimplePaginationOptions): JQuery;
    /**
     * The Vue Constructor
     * http://vuejs.org/api/index.html
     */
    constructor(options?: {});
    pagination(method: 'selectPage', pageNumber: number): void;
    pagination(method: 'prevPage'): void;
    pagination(method: 'nextPage'): void;
    pagination(method: 'getPagesCount'): number;
    pagination(method: 'getCurrentPage'): number;
    pagination(method: 'disable'): void;
    pagination(method: 'enable'): void;
    pagination(method: 'destroy'): void;
    pagination(method: 'redraw'): void;
    pagination(method: 'updateItems', items: number): void;
    pagination(method: string): any;
    pagination(method: string, value: any): any;
        /**
     * The Vue Constructor
     * http://vuejs.org/api/index.html
     */
    fun(a:{});
    paramAttributes:{}[];
    export function select(selector: string): Selection<any>;
    export var prototype: Selection<any>;
    attr(name: string, value: (datum: Datum, index: number) => Primitive): Update<Datum>;
    attr(obj: { [key: string]: Primitive | ((datum: Datum, index: number) => Primitive) }): Update<Datum>;
    property(obj: { [key: string]: any | ((datum: Datum, index: number) => any) }): Update<Datum>;
    call(func: (selection: Enter<Datum>, ...args: any[]) => any, ...args: any[]): Enter<Datum>;
}

interface ValueCallback {
    (newValue: {}, oldValue: {}): void;
  }

interface Group extends Array<EventTarget> {
            parentNode: EventTarget;
            [index: number]: Group;
        }

module test {
    export type Primitive = number | string | boolean; 
    tween(name: string, factory: () => (t: number) => any): Transition<Datum>;    
    style(obj: { [key: string]: Primitive | ((datum: Datum, index: number) => Primitive) }, priority?: string): Transition<Datum>;
      call(func: (transition: Transition<Datum>, ...args: any[]) => any, ...args: any[]): Transition<Datum>;
    export function ease(type: 'linear'): (t: number) => number;
    export function ease(type: 'linear-in'): (t: number) => number;
    export function mouse(container: EventTarget): [number, number];
}