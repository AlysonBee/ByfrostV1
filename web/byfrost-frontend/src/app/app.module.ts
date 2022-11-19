import { NgModule } from '@angular/core';
import { BrowserModule, By } from '@angular/platform-browser';

// import { FlexLayoutModule } from '@angular/flex-layout';
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';
import { ByfrostComponent } from './byfrost/byfrost.component';
import { AppRoutingModule } from './app-routing.module';
import { SanitizeHtmlPipe } from './byfrost/sanitizer.pipe';
import { ByfrostInfoComponent } from './byfrost/byfrost-info/byfrost-info.component';
import { FilesystemPipe } from './byfrost/filesystem.pipe';

@NgModule({
  declarations: [
    AppComponent,
    ByfrostComponent,
    SanitizeHtmlPipe,
    FilesystemPipe,
    ByfrostInfoComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
