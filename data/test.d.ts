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
 once(event: 'touchstart', fn: (event: interaction.InteractionEvent) => void, context?: any): EventEmitter;
     constructor(lineWidth: number, lineColor: number, lineAlpha: number, fillColor: number, fillAlpha: number, fill: boolean, shape: Circle | Rectangle | Ellipse | Polygon);
shape: Circle | Rectangle | Ellipse | Polygon;
shape: Circle | Rectangle | Ellipse | Polygon;
 once(event: 'touchstart', fn: (event: interaction.InteractionEvent) => void, context?: any): EventEmitter;
        once(event: string, fn: Function, context?: any): EventEmitter;
         shape: Circle | Rectangle | Ellipse | Polygon;
        type: number;

        clone(): GraphicsData;

}

 interface Group extends Array<EventTarget> {
            parentNode: EventTarget;
        }
interface Update<Datum> {}
class Update<Datum> {
     property(obj: { [key: string]: any | ((datum: Datum, index: number) => any) }): Update<Datum>;
    insert(name: (datum: Datum, index: number) => EventTarget, before: (datum: Datum, index: number) => EventTarget): Update<Datum>;
    data(): Datum[];
    data<NewDatum>(data: NewDatum[], key?: (datum: NewDatum, index: number) => string): Update<NewDatum>;
    export function touch(container: EventTarget, touches: TouchList, identifer: number): [number, number];
     export function touches(container: EventTarget, touches?: TouchList): Array<[number, number]>;
    export function min<T>(array: T[], accessor: (datum: T, index: number) => string): string;
    export function min<T, U extends Numeric>(array: T[], accessor: (datum: T, index: number) => U): U;

     export function extent<T extends Numeric>(array: Array<T | Primitive>): [T | Primitive, T | Primitive];

         export function deviation<T>(array: T[], accessor: (datum: T, index: number) => number): number;
    export var bisect: typeof bisectRight;

    export function mean<T>(array: T[], accessor: (datum: T, index: number) => number): number;
    export function extent<T>(array: T[], accessor: (datum: T, index: number) => number): [number, number];

        call(func: (transition: Transition<Datum>, ...args: any[]) => any, ...args: any[]): Transition<Datum>;

    export function bisector<T, U>(comparator: (a: T, b: U) => number): {
        left: (array: T[], x: U, lo?: number, hi?: number) => number;
        right: (array: T[], x: U, lo?: number, hi?: number) => number;
    }
    forEach(func: (value: string) => any): void;
    export function set(array: string[]): Set;
    export function merge<T>(arrays: T[][]): T[];
    rollup<U>(func: (values: T[]) => U): Nest<T>;
        map(array: T[]): { [key: string]: any };
}
