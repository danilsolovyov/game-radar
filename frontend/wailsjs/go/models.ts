export namespace color {
	
	export class RGBA {
	    R: number;
	    G: number;
	    B: number;
	    A: number;
	
	    static createFrom(source: any = {}) {
	        return new RGBA(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.R = source["R"];
	        this.G = source["G"];
	        this.B = source["B"];
	        this.A = source["A"];
	    }
	}

}

export namespace models {
	
	export class Blip {
	    angle: number;
	    distance: number;
	    intensity: number;
	
	    static createFrom(source: any = {}) {
	        return new Blip(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.angle = source["angle"];
	        this.distance = source["distance"];
	        this.intensity = source["intensity"];
	    }
	}
	export class DeviceSpeakers {
	    id: string;
	    name: string;
	    format_pcm: number;
	    rate: number;
	    channels: number;
	    default_period: number;
	    minimum_period: number;
	    latency: number;
	    is_default: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DeviceSpeakers(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.format_pcm = source["format_pcm"];
	        this.rate = source["rate"];
	        this.channels = source["channels"];
	        this.default_period = source["default_period"];
	        this.minimum_period = source["minimum_period"];
	        this.latency = source["latency"];
	        this.is_default = source["is_default"];
	    }
	}
	export class Theme {
	    name: string;
	    background_color: color.RGBA;
	    radar_color: color.RGBA;
	    border_opacity: number;
	    border_width: number;
	    section_base_opacity: number;
	    section_bright_opacity: number;
	    section_timeout: number;
	    section_count: number;
	    ring_count: number;
	    show_blips: boolean;
	    blip_opacity: number;
	    blip_timeout: number;
	    blip_size: number;
	    size: number;
	    pos_x: number;
	    pos_y: number;
	    intensity_multiplier: number;
	
	    static createFrom(source: any = {}) {
	        return new Theme(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.background_color = this.convertValues(source["background_color"], color.RGBA);
	        this.radar_color = this.convertValues(source["radar_color"], color.RGBA);
	        this.border_opacity = source["border_opacity"];
	        this.border_width = source["border_width"];
	        this.section_base_opacity = source["section_base_opacity"];
	        this.section_bright_opacity = source["section_bright_opacity"];
	        this.section_timeout = source["section_timeout"];
	        this.section_count = source["section_count"];
	        this.ring_count = source["ring_count"];
	        this.show_blips = source["show_blips"];
	        this.blip_opacity = source["blip_opacity"];
	        this.blip_timeout = source["blip_timeout"];
	        this.blip_size = source["blip_size"];
	        this.size = source["size"];
	        this.pos_x = source["pos_x"];
	        this.pos_y = source["pos_y"];
	        this.intensity_multiplier = source["intensity_multiplier"];
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

}

