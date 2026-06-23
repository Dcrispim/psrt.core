export namespace main {
	
	export class CompileSVGResult {
	    uri: string;
	    usedGoTextFallback: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CompileSVGResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uri = source["uri"];
	        this.usedGoTextFallback = source["usedGoTextFallback"];
	    }
	}

}

export namespace styleadapter {
	
	export class WebPreviewStyle {
	    container: Record<string, string>;
	    text: Record<string, string>;
	    hasStroke: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WebPreviewStyle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.container = source["container"];
	        this.text = source["text"];
	        this.hasStroke = source["hasStroke"];
	    }
	}

}

export namespace visualapp {
	
	export class MaskDetail {
	    index: number;
	    x: number;
	    y: number;
	    width: number;
	    height: number;
	    imageRef: string;
	    style: string;
	
	    static createFrom(source: any = {}) {
	        return new MaskDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.imageRef = source["imageRef"];
	        this.style = source["style"];
	    }
	}
	export class MaskPatch {
	    x?: number;
	    y?: number;
	    width?: number;
	    height?: number;
	    imageRef?: string;
	    styleSet?: Record<string, string>;
	    styleRemove?: string[];
	
	    static createFrom(source: any = {}) {
	        return new MaskPatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.imageRef = source["imageRef"];
	        this.styleSet = source["styleSet"];
	        this.styleRemove = source["styleRemove"];
	    }
	}
	export class OpenFileResult {
	    filePath: string;
	    document: string;
	
	    static createFrom(source: any = {}) {
	        return new OpenFileResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.document = source["document"];
	    }
	}
	export class PageDetail {
	    name: string;
	    imageUrl: string;
	    style: string;
	
	    static createFrom(source: any = {}) {
	        return new PageDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.imageUrl = source["imageUrl"];
	        this.style = source["style"];
	    }
	}
	export class PagePatch {
	    name?: string;
	    imageUrl?: string;
	    styleSet?: Record<string, string>;
	    styleRemove?: string[];
	
	    static createFrom(source: any = {}) {
	        return new PagePatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.imageUrl = source["imageUrl"];
	        this.styleSet = source["styleSet"];
	        this.styleRemove = source["styleRemove"];
	    }
	}
	export class PageSummary {
	    name: string;
	    imageUrl?: string;
	
	    static createFrom(source: any = {}) {
	        return new PageSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.imageUrl = source["imageUrl"];
	    }
	}
	export class TextDetail {
	    index: number;
	    x: number;
	    y: number;
	    width: number;
	    textSize: number;
	    content: string;
	    imageRef: string;
	    style: string;
	
	    static createFrom(source: any = {}) {
	        return new TextDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.textSize = source["textSize"];
	        this.content = source["content"];
	        this.imageRef = source["imageRef"];
	        this.style = source["style"];
	    }
	}
	export class TextPatch {
	    content?: string;
	    append: boolean;
	    x?: number;
	    y?: number;
	    width?: number;
	    textSize?: number;
	    imageRef?: string;
	    styleSet?: Record<string, string>;
	    styleRemove?: string[];
	
	    static createFrom(source: any = {}) {
	        return new TextPatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.append = source["append"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.textSize = source["textSize"];
	        this.imageRef = source["imageRef"];
	        this.styleSet = source["styleSet"];
	        this.styleRemove = source["styleRemove"];
	    }
	}
	export class UIState {
	    filePath: string;
	    activePage: string;
	    selectedIndex: number;
	    pages: PageSummary[];
	    page?: PageDetail;
	    texts?: TextDetail[];
	    masks?: MaskDetail[];
	    text?: TextDetail;
	    mask?: MaskDetail;
	    fonts: string[];
	    consts: Record<string, string>;
	    autoCompile: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UIState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.activePage = source["activePage"];
	        this.selectedIndex = source["selectedIndex"];
	        this.pages = this.convertValues(source["pages"], PageSummary);
	        this.page = this.convertValues(source["page"], PageDetail);
	        this.texts = this.convertValues(source["texts"], TextDetail);
	        this.masks = this.convertValues(source["masks"], MaskDetail);
	        this.text = this.convertValues(source["text"], TextDetail);
	        this.mask = this.convertValues(source["mask"], MaskDetail);
	        this.fonts = source["fonts"];
	        this.consts = source["consts"];
	        this.autoCompile = source["autoCompile"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VariantPSRT {
	    label: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new VariantPSRT(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.content = source["content"];
	    }
	}

}

