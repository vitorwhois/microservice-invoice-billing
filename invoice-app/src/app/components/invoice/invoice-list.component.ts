import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { BillingService } from '../../services/billing.service';
import { InventoryService } from '../../services/inventory.services';

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

  constructor(private billingService: BillingService, private inventoryService: InventoryService) { }

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
        this.inventoryService.getAllProducts();
      },
      error: (err) => alert('Erro ao adicionar item: ' + err.error)
    });
  }

  printInvoice(invoiceId: number): void {
    this.billingService.printInvoice(invoiceId)
      .subscribe({
        next: () => alert('Nota impressa com sucesso!'),
        error: (err) => {
          if (err.error) {
            const { step_reached, failed_reason, recovery } = err.error;

            let message = `Erro ao imprimir a nota fiscal.\n`;
            message += `Etapa com falha: ${step_reached}\n`;
            message += `Motivo: ${failed_reason}\n`;
            message += 'Erro: ' + err.error;

            if (recovery?.attempted) {
              message += `Recuperação: ${recovery.message} \n`;
              message += recovery.successful ? 'Recuperação bem-sucedida!' : 'A recuperação falhou!';
            }

            alert(message);
          } else {
            alert('Erro ao imprimir a nota fiscal. Tente novamente mais tarde.');
          }
        }
      });
  }



  openAddItemModal(invoiceId: number): void {
    this.selectedInvoiceId = invoiceId;
    this.showModal = true;
  }
}
