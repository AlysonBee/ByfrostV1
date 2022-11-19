import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { AppComponent } from "./app.component";
import { ByfrostComponent } from "./byfrost/byfrost.component";

const appRoutes: Routes = [
    { path: "", component: AppComponent},
    { path: "indexer", component: ByfrostComponent},
    { path: "**", redirectTo: "/indexer"}
]

@NgModule({
    imports: [
        RouterModule.forRoot(appRoutes)
    ],
    exports: [RouterModule]
})
export class AppRoutingModule{

}