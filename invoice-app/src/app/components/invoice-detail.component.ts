import { Component, Input } from '@angular/core';
import { BillingService } from '../services/billing.service';

@Component({
  selector: 'app-invoice-detail',
  templateUrl: './invoice-detail.component.html'
})
export class InvoiceDetailComponent {
  @Input() invoice: any;
  newItem = { productId: 0, quantity: 0 };

  constructor(private billingService: BillingService) { }

  addItem() {
    this.billingService.addItemToInvoice(
      this.invoice.ID,
      this.newItem.productId,
      this.newItem.quantity
    ).subscribe({
      next: () => this.loadInvoice(),
      error: (err) => alert('Erro: ' + err.error)
    });
  }

  printInvoice() {
    this.billingService.printInvoice(this.invoice.ID).subscribe({
      next: () => {
        alert('Nota impressa com sucesso!');
        this.loadInvoice();
      },
      error: (err) => alert('Erro: ' + err.error)
    });
  }

  private loadInvoice() {
    this.billingService.getInvoice(this.invoice.ID).subscribe(invoice => {
      this.invoice = invoice;
    });
  }
}
