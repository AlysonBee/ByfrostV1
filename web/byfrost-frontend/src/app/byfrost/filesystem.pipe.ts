import { Pipe, PipeTransform } from "@angular/core";
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';

@Pipe({
  name: 'filesystemPipe'
})
export class FilesystemPipe implements PipeTransform {

  constructor(private _sanitizer:DomSanitizer) {
  }

  transform(v:string):SafeHtml {
    let segments = v.split("/");
    let spaces = segments.length - 1;
    console.log("spaces length is ", spaces)
    if (spaces == 1) {
        return v
    }

    let prettier = "";
    while (spaces > 0) {
        prettier += "-"
        spaces--
    }
    prettier = prettier + "/" + segments[segments.length - 1]
    return prettier
   //  return this._sanitizer.bypassSecurityTrustHtml(v);
  }
}