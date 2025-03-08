import { TranslateModule } from '@ngx-translate/core';
import { Subject, takeUntil } from 'rxjs';

import { NgIf } from '@angular/common';
import { Component, inject, input, OnDestroy, OnInit } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatSidenavModule } from '@angular/material/sidenav';
import { RouterLink } from '@angular/router';

import { IconComponent } from '@app/components/icon';
import { Sidebar } from '@app/models';
import { SidebarService } from '@app/services/sidebar';
import { IS_MOBILE } from '@app/tokens';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss',
  imports: [
    MatSidenavModule,
    MatButtonModule,
    IconComponent,
    RouterLink,
    NgIf,
    TranslateModule,
  ],
})
export class SidebarComponent implements OnInit, OnDestroy {
  isCollapsed = true;

  items = input<Sidebar[]>([]);

  private sidebarService = inject(SidebarService);

  private onDestroy$ = new Subject();

  isMobile = inject(IS_MOBILE);

  ngOnInit(): void {
    this.sidebarService.isCollapsed$
      .pipe(takeUntil(this.onDestroy$))
      .subscribe(value => {
        this.isCollapsed = value;
      });
  }

  ngOnDestroy(): void {
    this.onDestroy$.unsubscribe();
  }
}
