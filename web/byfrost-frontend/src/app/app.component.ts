import { Component, ElementRef, OnInit, Renderer2} from '@angular/core';


@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'byfrost-frontend';

  ngOnInit() : void {
    console.log("here")
  } 
  constructor(private elementRef: ElementRef,
    private renderer2: Renderer2) {}
  ngAfterViewInit() {
      this.elementRef.nativeElement.ownerDocument
          .body.style.backgroundColor = 'black';
  }
}
