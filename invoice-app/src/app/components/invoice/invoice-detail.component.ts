import { Component, Input } from '@angular/core';
import { BillingService } from '../../services/billing.service';
import { InventoryService } from '../../services/inventory.services';

@Component({
  selector: 'app-invoice-detail',
  templateUrl: './invoice-detail.component.html'
})
export class InvoiceDetailComponent {
  @Input() invoice: any;
  newItem = { productId: 0, quantity: 0 };

  constructor(private billingService: BillingService, private inventoryService: InventoryService) { }

  addItem() {
    if (!this.newItem.productId || !this.newItem.quantity) {
      alert('Informe o ID do produto e a quantidade');
      return;
    }

    this.billingService.addItemToInvoice(this.invoice.ID, this.newItem.productId, this.newItem.quantity)
      .subscribe({
        next: () => {
          this.loadInvoice();
          this.updateProductList();
        },
        error: (err) => alert('Erro ao adicionar item: ' + err.error)
      });
  }

  private updateProductList() {
    this.inventoryService.getAllProducts().subscribe({
      next: (products) => console.log('Produtos atualizados:', products),
      error: (err) => console.error('Erro ao atualizar produtos:', err)
    });
  }

  printInvoice() {
    this.billingService.printInvoice(this.invoice.ID)
      .subscribe({
        next: () => {
          alert('Nota impressa com sucesso!');
          this.loadInvoice();
          this.updateProductList();
        },
        error: (err) => alert('Erro ao imprimir: ' + err.error)
      });
  }


  private loadInvoice() {
    this.billingService.getInvoice(this.invoice.ID)
      .subscribe(invoice => this.invoice = invoice);
  }
}
