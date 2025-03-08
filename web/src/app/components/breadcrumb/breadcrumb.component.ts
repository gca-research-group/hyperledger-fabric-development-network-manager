import { TranslateModule } from '@ngx-translate/core';

import { LowerCasePipe } from '@angular/common';
import { Component, input } from '@angular/core';
import { RouterModule } from '@angular/router';

import { Breadcrumb } from '@app/models';

@Component({
  selector: 'app-breadcrumb',
  templateUrl: './breadcrumb.component.html',
  styleUrl: './breadcrumb.component.scss',
  imports: [RouterModule, TranslateModule, LowerCasePipe],
})
export class BreadcrumbComponent {
  items = input<Breadcrumb[]>([]);
}
