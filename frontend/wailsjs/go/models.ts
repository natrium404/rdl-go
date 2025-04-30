export namespace models {
	
	export class Audio {
	    MimeType: string;
	    URL: string;
	    Codecs: string;
	
	    static createFrom(source: any = {}) {
	        return new Audio(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MimeType = source["MimeType"];
	        this.URL = source["URL"];
	        this.Codecs = source["Codecs"];
	    }
	}
	export class User {
	    Username: string;
	    Profile: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Username = source["Username"];
	        this.Profile = source["Profile"];
	    }
	}
	export class Video {
	    MimeType: string;
	    Codecs: string;
	    Quality: string;
	    URL: string;
	    Width: string;
	    Height: string;
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MimeType = source["MimeType"];
	        this.Codecs = source["Codecs"];
	        this.Quality = source["Quality"];
	        this.URL = source["URL"];
	        this.Width = source["Width"];
	        this.Height = source["Height"];
	    }
	}
	export class VideoData {
	    Caption: string;
	    Thumbnail: string;
	    User: User;
	    Videos: Video[];
	    Audio: Audio;
	    Reel: string;
	    Code: string;
	
	    static createFrom(source: any = {}) {
	        return new VideoData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Caption = source["Caption"];
	        this.Thumbnail = source["Thumbnail"];
	        this.User = this.convertValues(source["User"], User);
	        this.Videos = this.convertValues(source["Videos"], Video);
	        this.Audio = this.convertValues(source["Audio"], Audio);
	        this.Reel = source["Reel"];
	        this.Code = source["Code"];
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

export namespace scraper {
	
	export class Response {
	    Data: models.VideoData;
	    Success: boolean;
	    Message: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = this.convertValues(source["Data"], models.VideoData);
	        this.Success = source["Success"];
	        this.Message = source["Message"];
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

