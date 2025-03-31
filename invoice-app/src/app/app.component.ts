import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  selectedInvoice: any = null;

  setInvoice(invoice: any) {
    this.selectedInvoice = invoice;
  }
}
