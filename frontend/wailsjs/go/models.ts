export namespace play {
	
	export class PlaybackProgress {
	    positionMs: number;
	    durationMs: number;
	
	    static createFrom(source: any = {}) {
	        return new PlaybackProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.positionMs = source["positionMs"];
	        this.durationMs = source["durationMs"];
	    }
	}

}

