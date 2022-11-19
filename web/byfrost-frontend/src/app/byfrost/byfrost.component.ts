import { Component, OnInit, OnDestroy, Renderer2, ElementRef, ViewChild, QueryList, ViewChildren, AfterViewInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { ByfrostApiService } from './services/byfrost-api.service';
import { ByfrostTemplateService } from './services/byfrost-template.service';

import 'leader-line';
import { flatMap } from 'rxjs/operators';

declare let LeaderLine: any;

interface Coordinates {
  height: number
	width: number
	xcoord: number
	ycoord: number
}

interface ProgramConstruct {
  name: string;
  file: string;
  tokens: Token[];
  coord: Coordinates
}

interface Token {
  name: string;
  filepath: string;
  typ: string;
  line: number;
  spaces: number;
  tabs: number;
  label: string;
  id: string;
  styletag: string;
}

interface LineDrawing {
  start: any;
  end: any;
  theLine: any;
  isSet: boolean;
}

interface ScreenGrabLocations {
  functionName: string;
  top: number;
  left: number;
}

@Component({
  selector: 'app-byfrost',
  templateUrl: './byfrost.component.html',
  styleUrls: ['./byfrost.component.scss']
})
export class ByfrostComponent implements OnInit, OnDestroy, AfterViewInit{
  

  @ViewChildren("done") done: QueryList<any>;
  sidebarShow: boolean = false;

  codePayload: any;
  sub: Subscription; 
  progBody: string = "";
  preId: string = "/main";
  filePath: string = "";
  divList = new Map<string, string>();
  lineList = new Map<string, LineDrawing>();
  positionList = new Map<string, ScreenGrabLocations>();

  collapseLineList: string[] = [];
  clicked: string = "";
  entryPoint: string = "main";
  div: string = "";
  ProgContainer: ProgramConstruct;
  allLines: any[] = [];
  allLinearLines: any[] = [];  

  constructor(
    private byfrost: ByfrostApiService,
    private renderer2: Renderer2,
    private template: ByfrostTemplateService,
  ) { }

  private unlistener: () => void;

  // DEBUG
  getMethods(obj:any) {
    var result = [];
    for (var id in obj) {
      try {
        if (typeof(obj[id]) != "function") {
          result.push(id + ": " + obj[id].toString());
        }
    
      } catch (err) {
        result.push(id + ": inaccessible");
      }
    }
    return result;
  }


  leadLineHandler(lineID: string) {
    const startingElement = document.getElementById(lineID)
    const endingElement = document.getElementById(lineID.replace("-", "/"))
    if (endingElement == null) {
      return null
    }
    let newLine = new LeaderLine(startingElement, endingElement,
      {startPlugColor: '#99ff99', gradient: true,})  

    return newLine
  }

  collapseLinelistHandler(collapseList: string[]) {
    let newCollapseList: string[] = [];

    for (let i = 0; i < collapseList.length; i++) {
      let stringIndex = collapseList[i].length - 1;

      for (stringIndex; stringIndex > 0; stringIndex--) {
        if (collapseList[i][stringIndex] == '/') {
        
          let newString = collapseList[i].substring(0, stringIndex) +
            "-" + collapseList[i].substring(stringIndex + 1,
              collapseList[i].length)

      
          newCollapseList.push(newString);            
          break
        }
      }
    }
    console.log("newCollapseList is ", newCollapseList)
    return this.collapseLines(newCollapseList);
  }

  // DEBUG
  printDictionary(dictionary: any) {
    for (let key of dictionary.keys()) {
      console.log("key is ", key)
    }
  }

  collapseLines(collapseLineList: string[]) {
    console.log("lineList is ", this.lineList)
    for (let i = 0; i < collapseLineList.length; i++) {
      let theLine = this.lineList.get(collapseLineList[i]) 
    
      if (theLine) {
        theLine.theLine.remove(); 
        this.lineList.delete(collapseLineList[i])
      }
    }
  }

  private strncmp(str1: string, str2: string, n: number) {
    str1 = str1.substring(0, n);
    str2 = str2.substring(0, n);
    return ( ( str1 == str2 ) ? 0 : (( str1 > str2 ) ? 1 : -1 ));
  }

  drawLine(clicked: string):void  {
    const start = document.getElementById(clicked)
    const end = document.getElementById(clicked.replace("-", "/"))
    let drawLine: LineDrawing = {theLine: this.leadLineHandler(clicked), isSet: true,
      start:start, end:end}  
    this.lineList.set(clicked, drawLine);
  }

  repostionAll(): void {
    for (let [key, value] of this.lineList) {
      if (value.theLine) {
          value.theLine.position()
      }
    }
  }

  ngOnInit(): void {
    this.ProgContainer = this.newProgramConstruct(); 

    this.sub = this.byfrost.getCode(this.entryPoint)
      .subscribe((responseData) => {
        this.codePayload = responseData

        this.loadProgCotainer()
      });

    this.unlistener = this.renderer2.listen("document", "click", event => {
      if (event.currentTarget.activeElement.tagName.toLowerCase() === "button") {
        // console.log(`I am detecting mouseclick at ${event.currentTarget.activeElement}`);

        let idName = event.currentTarget.activeElement.id;
        this.clicked = idName;
        // If there's no ID, a link was clicked so get the name.
        if (idName.length == 0) {
          idName = event.currentTarget.activeElement.innerText; 
        }


        if (this.strncmp(idName, "visible", 7) != -1) {
          this.navigate(idName)
          return 
        } 

        this.byfrost.requestBlock(idName)
          .subscribe((responseData) => {
            
            let codeBody: any = responseData
            if (codeBody == null) {
              return 
            }
            this.codePayload = codeBody.Body
      
            if (codeBody.collapse == true) { 
              this.collapse(codeBody.collapselist)
              this.collapseLineList = codeBody.collapselist 
              console.log("collapsed list ", this.collapseLineList)
              this.collapseLinelistHandler(this.collapseLineList);

            } else {
              console.log("xcoord is ", this.ProgContainer.coord.ycoord)
              // window.scrollBy(600,this.ProgContainer.coord.ycoord);
              // window.scrollTo({
              //   top: 100,
              //   left: 10000,
              //   behavior: 'smooth'
              // })
              this.preId = idName.replace("-", "/")
              this.loadProgCotainer()
            }
          })
      }
    });
  }

  ngAfterViewInit(): void {
 
    this.done.changes.subscribe(t => {
      if (this.collapseLineList.length > 0) {
        this.collapseLineList = []
        return
      }
      if (this.clicked.length == 0) {
        return;
      }

      this.drawLine(this.clicked)
    });
    
  }


  collapse(targetList: string[]): void {
    for (let i = 0; i < targetList.length; i++) {
      let target = document.getElementById(targetList[i])
      if (target) {
        this.divList.delete(targetList[i])
        this.positionList.delete(targetList[i])
        target.remove();
        target.innerHTML = "";
      }
    }
  }


  ngOnDestroy(): void {
    this.sub.unsubscribe()
  }

  // this is for Ace editor
  //   ngAfterViewInit() {
  //     this.done.changes.subscribe(t => {
  //       ace.config.set("fontSize", "14px");
  //       const aceEditor = ace.edit("main"); 
  //       aceEditor.setReadOnly(true);
  //       ace.config.set('basePath', 'https://unpkg.com/ace-builds@1.4.12/src-noconflict');
  //       aceEditor.setTheme("ace/theme/monokai");
        

  //     })
  // }


  private loadProgCotainer() { 
    this.ProgContainer.tokens = this.codePayload.Tokens
    this.ProgContainer.name = this.codePayload.label
    this.ProgContainer.file = this.codePayload.path
    this.ProgContainer.coord = this.codePayload.Coord
    console.log(this.ProgContainer)
    this.filePath = this.ProgContainer.file; 
    this.buildProgramBody();
  }

  private newProgramConstruct():ProgramConstruct{
    return {
      name: "", file: "", tokens: [], 
      coord: {
        height: 0,
        width: 0,
        ycoord: 0,
        xcoord: 0
      }
    }
  }

  private repeatingChars(c: string, times: number) {
    let counter = 0;

    while (counter < times) {
      this.progBody += c;
      counter++;
    }
  }


  private newline(currline: number, linecount: number) {
      while (currline < linecount) {
        this.progBody += "\n";
        this.progBody += `<span style="color:white;">${String(currline)}</span>`;
        this.progBody += " ";
        currline++
      }
  }

  private buildProgramBody() {
    let currLine = this.ProgContainer.tokens[0].line-1
    // TRASH! FIX!

    this.progBody += "\n";
    this.progBody += `<span style="color:white;">${String(currLine)}</span>`;
    this.progBody += " ";
    let colorTag = "";
    
    currLine++
    for (let i = 0; i < this.ProgContainer.tokens.length; i++) {
      let tk = this.ProgContainer.tokens[i];

      if (tk.line > currLine) {
        this.newline(currLine, tk.line);
        currLine = tk.line;
      }
      
      if (tk.spaces > 0) {
        this.repeatingChars(" ", tk.spaces);
      }
      
      if (tk.tabs > 0) {
        this.repeatingChars("\t", tk.tabs);
      }




      if (tk.label == "DECL" && this.preId.length == 0) { 
        this.preId = tk.id;
      }
      if (tk.label == "PARAM_TYPE") {
        console.log("LABEL IS ", tk.label)
      }

      colorTag = this.template.setColour(tk.styletag, tk.name)
      
      if (tk.label == "FUNCTION" || tk.label == "DEREF" || tk.label == "PARAM_TYPE") {
        this.progBody = this.progBody + this.template.getLabelLink("label", colorTag, tk.id);
      } else {
        this.progBody += colorTag;  
      }
    }

    this.div = this.div + this.template.getPreBody(
        this.filePath, this.progBody, this.preId,
        this.ProgContainer.coord.ycoord+25,
        this.ProgContainer.coord.xcoord+30
    );
    this.divList.set(this.preId, this.div)
    this.div = ""
    this.progBody = ""
  //  this.preId = ""

    window.scrollTo({
      top: this.ProgContainer.coord.ycoord-100,
      left: this.ProgContainer.coord.xcoord-100,
      behavior: 'smooth' 
    })

    this.saveVisiblePositions(
      this.preId, this.preId,
      this.ProgContainer.coord.ycoord-100,
      this.ProgContainer.coord.xcoord-100
    )
    this.preId = ""
  }

 private navigate(functionName: string) {
   let functionId = functionName.split(":")
   if (functionId.length == 2) {
     let position = this.positionList.get(functionId[1])
     window.scrollTo({
      top: position?.top,
      left: position?.left,
      behavior: 'smooth' 
    })
  }


}


 private saveVisiblePositions(functionName: string, functionId: string, top: number, left: number) {
    this.positionList.set(
      functionId, 
      {
        functionName: functionName,
        top: top,
        left: left
      }
    )
    
  }
}
