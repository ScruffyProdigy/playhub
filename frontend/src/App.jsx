import CompleteSignInPage from './components/auth/CompleteSignInPage'
import ReturnPage from './components/auth/ReturnPage'
import AuthPanel from './components/auth/AuthPanel'
import ActiveQueueBanner from './components/games/ActiveQueueBanner'
import GameLobby from './components/games/GameLobby'
import CreateRoomPanel from './components/rooms/CreateRoomPanel'
import { ActiveRoomProvider, useActiveRoom } from './components/rooms/ActiveRoomProvider'
import AppDock from './components/rooms/AppDock'
import RoomPanel from './components/rooms/RoomPanel'
import RoomSheet from './components/rooms/RoomSheet'
import { AuthProvider, useAuth } from './components/auth/AuthProvider'
import { useActiveQueue } from './components/games/useActiveQueue'
import { useActiveTableSeat } from './components/games/useActiveTableSeat'
import { APP_NAME, APP_TAGLINE } from './lib/brand'
import { parseRoomInviteCode } from './lib/rooms'
import { MOBILE_ROOM_QUERY, useMediaQuery } from './lib/useMediaQuery'
import { usePathname } from './lib/usePathname'
import { useEffect } from 'react'
import './App.css'

function CatalogPage() {
  const { user, loading: authLoading } = useAuth()
  const { activeQueue, busy: queueBusy, refresh: refreshQueue, handleLeave: leaveQueue } = useActiveQueue()
  const { activeTableSeat, busy: tableBusy, refresh: refreshTable, handleLeave: leaveTableSeat } =
    useActiveTableSeat()

  const bannerBusy = queueBusy || tableBusy
  const onLeaveIntent = activeQueue?.queueId ? leaveQueue : leaveTableSeat

  return (
    <main className="app-shell app-shell--catalog">
      {!authLoading && user ? (
        <ActiveQueueBanner
          activeQueue={activeQueue}
          activeTableSeat={activeTableSeat}
          busy={bannerBusy}
          onLeave={onLeaveIntent}
        />
      ) : null}

      <header className="app-header">
        <h1>{APP_NAME}</h1>
        <p className="tagline">{APP_TAGLINE}</p>
      </header>

      <AuthPanel />
      {!authLoading && user ? <CreateRoomPanel /> : null}
      <GameLobby
        activeQueue={activeQueue}
        activeTableSeat={activeTableSeat}
        onQueueChange={refreshQueue}
        onTableChange={refreshTable}
      />
    </main>
  )
}

function MainLayout() {
  const pathname = usePathname()
  const inviteCode = parseRoomInviteCode(pathname)
  const { user } = useAuth()
  const { room, roomOpen, dismissRoom, openRoom, unreadCount, hasRoomMembership } = useActiveRoom()
  const isMobile = useMediaQuery(MOBILE_ROOM_QUERY)
  const inRoomContext = Boolean(room || inviteCode || hasRoomMembership)
  const showDesktopRoom = Boolean(room && user && !isMobile && roomOpen)
  const showMobileSheet = Boolean(user && isMobile && roomOpen)
  const showDock = Boolean(user && isMobile && inRoomContext)

  useEffect(() => {
    const root = document.getElementById('root')
    root?.classList.toggle('app-root--split', showDesktopRoom)
    root?.classList.toggle('app-root--dock', showDock)
    return () => {
      root?.classList.remove('app-root--split', 'app-root--dock')
    }
  }, [showDesktopRoom, showDock])

  return (
    <>
      <div className={`app-layout ${showDesktopRoom ? 'app-layout--split' : ''}`}>
        <div className="app-layout__catalog">
          <CatalogPage />
        </div>
        {showDesktopRoom ? (
          <aside className="app-layout__room panel-card" aria-label="Room chat">
            <RoomPanel compact />
          </aside>
        ) : null}
      </div>

      {showMobileSheet ? (
        <RoomSheet open={roomOpen} onDismiss={dismissRoom}>
          <RoomPanel compact />
        </RoomSheet>
      ) : null}

      {showDock ? (
        <AppDock
          unreadCount={unreadCount}
          roomOpen={roomOpen}
          onOpenCatalog={dismissRoom}
          onOpenRoom={openRoom}
        />
      ) : null}
    </>
  )
}

function MainShell() {
  const pathname = usePathname()
  const inviteCode = parseRoomInviteCode(pathname)

  return (
    <ActiveRoomProvider pendingInviteCode={inviteCode}>
      <MainLayout />
    </ActiveRoomProvider>
  )
}

function App() {
  const pathname = usePathname()
  const isMainRoute = pathname === '/' || parseRoomInviteCode(pathname)

  return (
    <AuthProvider>
      {pathname.startsWith('/auth/complete') ? (
        <CompleteSignInPage />
      ) : pathname.startsWith('/return') ? (
        <ReturnPage />
      ) : isMainRoute ? (
        <MainShell />
      ) : (
        <main className="app-shell auth-page">
          <h1>Page not found</h1>
          <a className="auth-link" href="/">
            Back to {APP_NAME}
          </a>
        </main>
      )}
    </AuthProvider>
  )
}

export default App
