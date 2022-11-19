import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ByfrostTemplateService {
  top = 0

  // This is more documentation so that i remember the values.
  buttonStyle = `background: none;
	color: white;
	border: 1px solid #ff66ff;
	padding: 0;
	font: inherit;
	cursor: pointer;
	outline: inherit;`

  constructor() { }
 
  private getButtonStyle(color: string) {
    return `background: none;
    color: ${color};
    border: none;
    padding: 0;
    font: inherit;
    cursor: pointer;
    outline: inherit;` 
  }
 
  setColour(color: string, html: string) {
    return `<span style="color:${color}">${html}</span>`
  }

  getLabelLink(className: string, linkName: string, id: string) {
    return `<button id="${id}" class="${className}" style=\"${this.getButtonStyle("inherit")}\ color="white";>${linkName}</button>`
  }

  getDerefLink(className: string, linkName: string, id: string) {
    return `<button id="${id}" class="${className}" style=\"${this.getButtonStyle("white")}\">${linkName}</button>`
  }

  getPreBody(filePath: string, content: string, preId: string, y: number, x: number) {
    console.log("preId is ", preId)
    return `<div id="${preId}" style="padding:30px; margin-top:${y}px; margin-left:${x}px; position:absolute; overflow:visible">
              <button style="${this.buttonStyle}">
               location:${filePath}
              </button>
              <div style="display:flex;">
                <pre style="background-color:black;">${content}</pre>
              </div>
            </div>`
  }

}
