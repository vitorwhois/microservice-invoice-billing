import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';
import { InvoiceListComponent } from './components/invoice-list.component';
import { InvoiceDetailComponent } from './components/invoice-detail.component';
import { ProductListComponent } from './components/product-list.component';


@NgModule({
  declarations: [
    AppComponent,
    InvoiceListComponent,
    InvoiceDetailComponent,
    ProductListComponent

  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
