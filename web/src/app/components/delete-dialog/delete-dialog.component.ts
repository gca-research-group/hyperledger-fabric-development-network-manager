import { TranslateModule } from '@ngx-translate/core';

import { Component, inject } from '@angular/core';
import { MatButton } from '@angular/material/button';
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle,
} from '@angular/material/dialog';

@Component({
  selector: 'app-delete-dialog',
  templateUrl: './delete-dialog.component.html',
  styleUrl: './delete-dialog.component.scss',
  imports: [
    MatDialogTitle,
    MatDialogContent,
    MatDialogActions,
    MatButton,

    TranslateModule,
  ],
})
export class DeleteDialogComponent<T> {
  readonly dialogRef = inject(MatDialogRef<DeleteDialogComponent<T>>);

  readonly data = inject<T>(MAT_DIALOG_DATA);

  cancel(): void {
    this.dialogRef.close();
  }

  ok() {
    this.dialogRef.close(this.data);
  }
}
