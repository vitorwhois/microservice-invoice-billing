import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { BillingService } from '../services/billing.service';

@Component({
  selector: 'app-invoice-list',
  templateUrl: './invoice-list.component.html'
})
export class InvoiceListComponent implements OnInit {
  invoices: any[] = [];
  newInvoiceNumber: string = '';
  showModal: boolean = false;
  selectedInvoiceId: number = 0;
  productId: number = 0;
  quantity: number = 0;



  @Output() invoiceSelected = new EventEmitter<any>();

  constructor(private billingService: BillingService) { }

  ngOnInit(): void {
    this.loadInvoices();
  }

  loadInvoices(): void {
    this.billingService.getInvoices().subscribe({
      next: (data) => this.invoices = data,
      error: (err) => alert('Erro ao carregar notas: ' + err.error)
    });
  }

  createInvoice(): void {
    if (!this.newInvoiceNumber) {
      alert('Informe um número para a nota fiscal');
      return;
    }
    this.billingService.createInvoice(this.newInvoiceNumber).subscribe({
      next: (invoice) => {
        alert('Nota fiscal criada!');
        this.newInvoiceNumber = '';
        this.loadInvoices();
      },
      error: (err) => alert('Erro ao criar nota: ' + err.error)
    });
  }



  selectInvoice(invoice: any): void {
    this.invoiceSelected.emit(invoice);
  }

  addItemToInvoice(invoiceId: number, productId: number, quantity: number): void {
    if (invoiceId === null || productId === null || quantity === null) {
      alert("Preencha todos os campos!");
      return;
    }

    this.billingService.addItemToInvoice(invoiceId, productId, quantity).subscribe({
      next: () => {
        alert('Item adicionado com sucesso!');
        this.showModal = false;
        this.loadInvoices();
      },
      error: (err) => alert('Erro ao adicionar item: ' + err.error)
    });
  }

  printInvoice(invoiceId: number): void {
    this.billingService.printInvoice(invoiceId)
      .subscribe({
        next: () => alert('Nota impressa com sucesso!'),
        error: (err) => alert('Erro ao imprimir: ' + err.error)
      });
  }

  openAddItemModal(invoiceId: number): void {
    this.selectedInvoiceId = invoiceId;
    this.showModal = true;
  }
  formatItems(items: any[]): string {
    return items.map(item => `${item.Name} (Qtd: ${item.Quantity}, Preço: ${item.Price})`).join(', ');
  }
}
