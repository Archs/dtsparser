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
      map(array: T[]): { [key: string]: any };
        map(array: T[], mapType: typeof d3.map): Map<any>;
        entries(array: T[]): { key: string; values: any }[];
}

interface Transform {
        rotate: number;
        translate: [number, number];
        skew: number;
        scale: [number, number];
        toString(): string;
        new (r: number, g: number, b: number): Rgb;
        export function requote(string: string): string;

        export var rgb: {
            new (r: number, g: number, b: number): Rgb;
            new (color: string): Rgb;
            (r: number, g: number, b: number): Rgb;
            (color: string): Rgb;
        };
        interpolate(): string | ((points: Array<[number, number]>) => string);
        interpolate(interpolate: "linear"): Line<T>;
        interpolate(interpolate: (points: Array<[number, number]>) => string): Radial<T>;
        export function radial(): Radial<Link<Node>, Node>;
        projection(projection: (d: Node, i: number) => [number, number]): Radial<Link, Node>;
         (url: string, callback: (error: any, rows: { [key: string]: string }[]) => void): DsvXhr<{ [key: string]: string }>;
        (url: string): DsvXhr<{ [key: string]: string }>;
        <T>(url: string, accessor: (row: { [key: string]: string }) => T, callback: (rows: T[]) => void): DsvXhr<T>;

        row<U>(accessor: (row: { [key: string]: string }) => U): DsvXhr<U>;
    }

interface Locale {
        numberFormat(specifier: string): (n: number) => string;
        timeFormat: {
            (specifier: string): time.Format;
            utc(specifier: string): time.Format;
        }
    }

interface Bundle<T extends bundle.Node> {
            (links: bundle.Link<T>[]): T[][];
            order(order: (data: Array<[number, number]>) => number[]): Stack<Series, Value>;
            links(nodes: T[]): tree.Link<T>[];
            value(): (datum: T, index: number) => number;
            value(value: (datum: T, index: number) => number): Tree<T>;
        }

interface Link<T extends Node> {
                source: T;
                target: T;
            }

module asdfa {
    type Padding = number | [number, number, number, number];
    visit(callback: (node: Node<T>, x1: number, y1: number, x2: number, y2: number) => boolean | void): void;
    clip(subject: Array<[number, number]>): Array<[number, number]>;
}

interface TouchList { }

declare module 'd3' {
    export = d3;
     _init(options: {}): void;
    _cleanup(): void;
    // static require(module:string) : void;
}

 interface FilterCallback {
    (value:{},begin?:{},end?:{}): {};
  }

import Vue = vuejs.Vue;

declare module "vue" {
    import vue = vuejs.Vue;
    export = vue;
}

interface KnockoutSubscribable<T> extends KnockoutSubscribableFunctions<T> {
    subscribe(callback: (newValue: T) => void, target?: any, event?: string): KnockoutSubscription;
    subscribe<TEvent>(callback: (newValue: TEvent) => void, target: any, event: string): KnockoutSubscription;
    extend(requestedExtenders: { [key: string]: any; }): KnockoutSubscribable<T>;
    getSubscriptionsCount(): number;
    compareArrays<T>(a: T[], b: T[]): Array<KnockoutArrayChange<T>>;
}

 interface Loader {
        getConfig? (componentName: string, callback: (result: ComponentConfig) => void): void;
        loadComponent? (componentName: string, config: ComponentConfig, callback: (result: Definition) => void): void;
        loadTemplate? (componentName: string, templateConfig: any, callback: (result: Node[]) => void): void;
        loadViewModel? (componentName: string, viewModelConfig: any, callback: (result: any) => void): void;
        suppressLoaderExceptions?: boolean;
    }

interface Definition {
        template: Node[];
        createViewModel? (params: any, options: { element: Node; }): any;
        get(componentName: string, callback: (definition: KnockoutComponentTypes.Definition) => void): void;
        loaders: KnockoutComponentTypes.Loader[];
    getComponentNameForNode(node: Node): string;

    }

declare var ko: KnockoutStatic;

declare module "knockout" {
    export = ko;
    variables: {
                clang: number;
                host_arch: string;
                node_install_npm: boolean;
                node_install_waf: boolean;
                node_prefix: string;
                node_shared_openssl: boolean;
                node_shared_v8: boolean;
                node_shared_zlib: boolean;
                node_use_dtrace: boolean;
                node_use_etw: boolean;
                node_use_openssl: boolean;
                target_arch: string;
                v8_no_strict_aliasing: number;
                v8_use_snapshot: boolean;
                visibility: string;
            };
     memoryUsage(): { rss: number; heapTotal: number; heapUsed: number; };
     send?(message: any, sendHandle?: any): void;
     Float32Array: typeof Float32Array;
    Float64Array: typeof Float64Array;
    Function: typeof Function;
    setTimeout: (callback: (...args: any[]) => void, ms: number, ...args: any[]) => NodeJS.Timer;
    undefined: typeof undefined;
        unescape: (str: string) => string;
        gc: () => void;
    export function parse(str: string, sep?: string, eq?: string, options?: { maxKeys?: number; }): any;
    emit(event: string, ...args: any[]): boolean;
    import events = require("events");
    import net = require("net");
    import stream = require("stream");
        /**
    * QUnit has a bunch of internal configuration defaults, some of which are 
    * useful to override. Check the description for each option for details.
    */
    config: Config;
        /**
    * QUnit has a bunch of internal configuration defaults, some of which are 
    * useful to override. Check the description for each option for details.
    */
    config: Config;
}

/**
* @param title Title of unit being tested
* @param test Function to close over assertions
*/
declare function test(title: string, test: (assert?: QUnitAssert) => any): any;

declare function notPropEqual(actual: any, expected: any, message?: string): any;

declare function propEqual(actual: any, expected: any, message?: string): any;

// https://github.com/jquery/qunit/blob/master/qunit/qunit.js#L1568
declare function equiv(a: any, b: any): any;

// https://github.com/jquery/qunit/blob/master/qunit/qunit.js#L661
declare var raises: any;

/* QUNIT */
declare var QUnit: QUnitStatic;

interface DoneCallbackObject {
    /**
    * The number of failed assertions
    */
    failed: number;

    /**
    * The number of passed assertions
    */
    passed: number;

    /**
    * The total number of assertions
    */
    total: number;

    /**
    * The time in milliseconds it took tests to run from start to finish.
    */
    runtime: number;
    
    /**
    * Alias of throws.
    * 
    * In very few environments, like Closure Compiler, throws is considered a reserved word
    * and will cause an error. For that case, an alias is bundled called raises. It has the
    * same signature and behaviour, just a different name.
    * 
    * @param block Function to execute
    * @param expected Error Object to compare
    * @param message A short description of the assertion
    */
    raises(block: () => any, expected: any, message?: string): any;
}