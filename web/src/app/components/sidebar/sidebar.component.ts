import { Component, inject, input, OnDestroy, OnInit } from '@angular/core';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatButtonModule } from '@angular/material/button';
import { IconComponent } from '@app/components/icon';
import { NgIf } from '@angular/common';
import { RouterLink } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { SidebarService } from '@app/services/sidebar';
import { Subject, takeUntil } from 'rxjs';
import { Sidebar } from '@app/models';
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
